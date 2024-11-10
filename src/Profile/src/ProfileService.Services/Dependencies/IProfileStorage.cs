using ProfileService.Services.Entities;

namespace ProfileService.Services.Dependencies;

public interface IProfileStorage
{
    Task<ProfileEntity> GetProfile(long profileId, CancellationToken cancellationToken);
    Task<long> CreateProfile(CreateProfileEntity profile, CancellationToken cancellationToken);
    Task UpdateProfile(ProfileEntity profile, CancellationToken cancellationToken);
    Task<IList<ProfileEntity>> GetProfiles(IList<long> profileIds, CancellationToken cancellationToken);
}