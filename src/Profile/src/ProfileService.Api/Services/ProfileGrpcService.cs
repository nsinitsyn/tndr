using Grpc.Core;
using Microsoft.AspNetCore.Authorization;
using ProfileService.Api.Authentication;
using ProfileService.Services.Entities;
using TinderApiV1;
using Service = ProfileService.Services.ProfileService;
using Gender = ProfileService.Services.Entities.Gender;
using ApiGender = TinderApiV1.Gender;

namespace ProfileService.Api.Services;

[Authorize]
public class ProfileGrpcService : TinderApiV1.ProfileService.ProfileServiceBase
{
    private readonly ILogger<ProfileGrpcService> _logger;
    private readonly IUserProfileProvider _userProfileProvider;
    private readonly Service _profileService;

    public ProfileGrpcService(
        ILogger<ProfileGrpcService> logger,
        IUserProfileProvider userProfileProvider,
        Service profileService)
    {
        _logger = logger;
        _userProfileProvider = userProfileProvider;
        _profileService = profileService;
    }

    // grpcurl -H 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQcm9maWxlSWQiOiIxMCIsImV4cCI6MTc2MjcxOTY2NSwiaXNzIjoiQXV0aFNlcnZlciIsImF1ZCI6IkF1dGhDbGllbnQifQ.pDWIoPzTE9Q_ccKgC11CMiczkKx52dYikYZEC6qvAbU' -plaintext -d '{}' 172.24.48.1:2340 tinder.ProfileService/GetMyProfile
    public override async Task<GetMyProfileResponse> GetMyProfile(GetMyProfileRequest request, ServerCallContext context)
    {
        var profile = await _profileService.GetProfile(_userProfileProvider.ProfileId, context.CancellationToken);

        var result = new GetMyProfileResponse
        {
            Profile = new ProfileDto
            {
                Gender = (ApiGender)profile.Gender,
                Age = profile.Age,
                Name = profile.Name,
                Description = profile.Description,
            }
        };
        result.Profile.PhotoUrls.Add(profile.Photos);

        return result;
    }

    // grpcurl -H 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQcm9maWxlSWQiOiIxMCIsImV4cCI6MTc2MjcxOTY2NSwiaXNzIjoiQXV0aFNlcnZlciIsImF1ZCI6IkF1dGhDbGllbnQifQ.pDWIoPzTE9Q_ccKgC11CMiczkKx52dYikYZEC6qvAbU' -plaintext -d '{"profile":{"age":32, "name":"Alex", "gender":"M", "description":"My profile", "photo_urls":["123", "456"]}}' 172.24.48.1:2340 tinder.ProfileService/CreateProfile
    public override async Task<CreateProfileResponse> CreateProfile(CreateProfileRequest request, ServerCallContext context)
    {
        var profileId = await _profileService.CreateProfile(new CreateProfileEntity
        {
            Gender = (Gender)request.Profile.Gender,
            Age = request.Profile.Age,
            Name = request.Profile.Name,
            Description = request.Profile.Description,
            Photos = request.Profile.PhotoUrls.ToList()
        }, context.CancellationToken);

        return new CreateProfileResponse
        {
            ProfileId = profileId
        };
    }

    // grpcurl -H 'authorization: Bearer <jwt_token>' -plaintext -d '{"Profile":null}' 172.24.48.1:2340 tinder.ProfileService/UpdateProfile
    // grpcurl -H 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQcm9maWxlSWQiOiIxMCIsImV4cCI6MTc2MjcxOTY2NSwiaXNzIjoiQXV0aFNlcnZlciIsImF1ZCI6IkF1dGhDbGllbnQifQ.pDWIoPzTE9Q_ccKgC11CMiczkKx52dYikYZEC6qvAbU' -plaintext -d '{"profile":{"age":28, "name":"Max", "description":"fff", "photo_urls":["123"]}}' 172.24.48.1:2340 tinder.ProfileService/UpdateProfile
    public override async Task<UpdateProfileResponse> UpdateProfile(UpdateProfileRequest request, ServerCallContext context)
    {
        await _profileService.UpdateProfile(new ProfileEntity
        {
            ProfileId = _userProfileProvider.ProfileId,
            Gender = (Gender)request.Profile.Gender,
            Age = request.Profile.Age,
            Name = request.Profile.Name,
            Description = request.Profile.Description,
            Photos = request.Profile.PhotoUrls.ToArray()
        }, context.CancellationToken);

        return new UpdateProfileResponse();
    }

    // Administrator role: grpcurl -H 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJSb2xlIjoiQWRtaW5pc3RyYXRvciIsImV4cCI6MTc2Mjc4MjM5OCwiaXNzIjoiQXV0aFNlcnZlciIsImF1ZCI6IkF1dGhDbGllbnQifQ.rVqdEV044q5b4QtXu4fJsGsyh1bFk8zp_EiPFwzh5LE' -plaintext -d '{"profile_ids":[10, 17, 18]}' 172.24.48.1:2340 tinder.ProfileService/GetProfiles
    [Authorize("AdministratorOrMatchService")]
    public override async Task<GetProfilesResponse> GetProfiles(GetProfilesRequest request, ServerCallContext context)
    {
        var profiles = await _profileService.GetProfiles(request.ProfileIds.ToList(), context.CancellationToken);

        var result = new GetProfilesResponse();

        result.Profiles.Add(profiles.Select(x =>
        {
            var r = new ProfileGetDto
            {
                ProfileId = x.ProfileId,
                Name = x.Name,
                Description = x.Description,
            };
            r.PhotoUrls.Add(x.Photos);
            return r;
        }));
        
        return result;
    }
}