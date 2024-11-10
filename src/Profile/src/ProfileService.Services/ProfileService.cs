using Microsoft.Extensions.Logging;
using ProfileService.Services.Dependencies;
using ProfileService.Services.Entities;
using ProfileService.Services.Exceptions;

namespace ProfileService.Services;

public class ProfileService
{
    private readonly ILogger<ProfileService> _logger;
    private readonly IProfileStorage _profileStorage;

    public ProfileService(ILogger<ProfileService> logger, IProfileStorage profileStorage)
    {
        _logger = logger;
        _profileStorage = profileStorage;
    }
    
    public async Task<ProfileEntity> GetProfile(long profileId, CancellationToken cancellationToken)
    {
        return await _profileStorage.GetProfile(profileId, cancellationToken);
    }
    
    public async Task<long> CreateProfile(CreateProfileEntity profile, CancellationToken cancellationToken)
    {
        return await _profileStorage.CreateProfile(profile, cancellationToken);
    }
    
    public async Task UpdateProfile(ProfileEntity profile, CancellationToken cancellationToken)
    {
        await _profileStorage.UpdateProfile(profile, cancellationToken);
    }
    
    public async Task<IList<ProfileEntity>> GetProfiles(IList<long> profileIds, CancellationToken cancellationToken)
    {
        var result = await _profileStorage.GetProfiles(profileIds, cancellationToken);

        if (result.Count == profileIds.Count)
        {
            return result;
        }
        
        _logger.LogError(
            "Unexpected behavior. Inconsistent result in GetProfiles method. " +
            "Requested count: {profileIdsCount}. " +
            "Resulted count: {resultCount}. " +
            "Requested data: {@profileIds}. " +
            "Result: {@result}. ",
            profileIds.Count, result.Count, profileIds, result);
            
        throw new NotLoggableException();

    }
}