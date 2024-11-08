using Grpc.Core;
using Microsoft.AspNetCore.Authorization;
using ProfileService.Api.Authentication;
using ProfileService.Domain;
using TinderApiV1;
using Service = ProfileService.Services.ProfileService;

namespace ProfileService.Api.GrpcServices;

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

    // grpcurl -H 'authorization: Bearer <jwt_token>' -plaintext -d '{"Profile":null}' 172.24.48.1:2340 tinder.ProfileService/UpdateProfile
    public override async Task<UpdateProfileResponse> UpdateProfile(UpdateProfileRequest request, ServerCallContext context)
    {
        await _profileService.UpdateProfile(new ProfileEntity
        {
            ProfileId = _userProfileProvider.ProfileId,
            Sex = request.Profile.Sex,
            Age = request.Profile.Age,
            Name = request.Profile.Name,
            Description = request.Profile.Description,
            Photos = request.Profile.PhotoUrls.ToList()
        }, context.CancellationToken);

        return new UpdateProfileResponse();
    }

    public override Task<GetProfilesResponse> GetProfiles(GetProfilesRequest request, ServerCallContext context)
    {
        return Task.FromResult(new GetProfilesResponse
        {
            UnexpectedServerError = new UnexpectedServerError()
        });
    }
}