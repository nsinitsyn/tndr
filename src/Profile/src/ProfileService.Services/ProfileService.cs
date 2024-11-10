using ProfileService.Domain;
using ProfileService.Services.Dependencies;

namespace ProfileService.Services;

public class ProfileService
{
    private readonly IProfileStorage _profileStorage;

    public ProfileService(IProfileStorage profileStorage)
    {
        _profileStorage = profileStorage;
    }
    
    public Task<long> AddProfile(ProfileEntity profile, CancellationToken cancellationToken)
    {
        throw new NotImplementedException();
    }
    
    public async Task UpdateProfile(ProfileEntity profile, CancellationToken cancellationToken)
    {
        await _profileStorage.UpdateProfile(profile);
    }
    
    public Task GetProfiles(IList<long> profileId, CancellationToken cancellationToken)
    {
        throw new NotImplementedException();
    }
}