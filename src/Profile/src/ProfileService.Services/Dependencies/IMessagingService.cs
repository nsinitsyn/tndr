using ProfileService.Services.Entities.Messaging;

namespace ProfileService.Services.Dependencies;

public interface IMessagingService
{
    Task Publish(ProfileUpdatedMessage message);
}