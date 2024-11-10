using ProfileService.Services.Entities;

namespace ProfileService.Services.Dependencies;

public interface IProfileOutboxStorage
{
    Task<IList<ProfileOutboxEntity>> GetProfileOutbox(int limit);
    Task ClearProfileOutbox(List<long> ids);
}