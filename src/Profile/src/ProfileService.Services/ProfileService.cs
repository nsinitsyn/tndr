using ProfileService.Domain;
using ProfileService.Services.Entities.Messaging;

namespace ProfileService.Services;

public interface IProfileRepository
{
    Task AddProfile(ProfileEntity profile);
    Task UpdateProfile(ProfileEntity profile);
    Task<ProfileEntity> GetProfiles();
}

public interface IProfileQueueNotifier
{
    Task SendProfileUpdatedMessage(ProfileUpdatedQueueMessage profileUpdatedMessage);
}

public class ProfileService
{
    private readonly IProfileRepository _profileRepository;
    private readonly IProfileQueueNotifier _profileQueueNotifier;

    public ProfileService(IProfileRepository profileRepository, IProfileQueueNotifier profileQueueNotifier)
    {
        _profileRepository = profileRepository;
        _profileQueueNotifier = profileQueueNotifier;
    }
    
    public Task<long> AddProfile(ProfileEntity profile, CancellationToken cancellationToken)
    {
        throw new NotImplementedException();
    }
    
    public Task UpdateProfile(ProfileEntity profile, CancellationToken cancellationToken)
    {
        throw new NotImplementedException();
    }
    
    public Task GetProfiles(IList<long> profileId, CancellationToken cancellationToken)
    {
        throw new NotImplementedException();
    }
}