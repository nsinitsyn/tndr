using ProfileService.Services.Entities;

namespace ProfileService.Services.Dependencies;

public interface IProfileOutboxStorage
{
    Task<IList<ProfileOutboxEntity>> GetProfileOutbox(int limit, CancellationToken cancellationToken);
    Task ClearProfileOutbox(List<long> ids, CancellationToken cancellationToken);
}