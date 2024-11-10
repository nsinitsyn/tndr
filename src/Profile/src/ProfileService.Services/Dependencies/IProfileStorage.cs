using ProfileService.Domain;

namespace ProfileService.Services.Dependencies;

public interface IProfileStorage
{
    Task AddProfile(ProfileEntity profile);
    Task UpdateProfile(ProfileEntity profile);
    Task<ProfileEntity> GetProfiles();
}