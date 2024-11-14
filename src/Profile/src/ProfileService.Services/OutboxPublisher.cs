using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using ProfileService.Services.Dependencies;
using ProfileService.Services.Entities;
using ProfileService.Services.Entities.Messaging;

namespace ProfileService.Services;

// Outbox pattern implementation with guarantee is send at least once.
// ProfileUpdate messages are idempotent currently.
public class OutboxPublisher : BackgroundService
{
    private readonly ILogger<OutboxPublisher> _logger;
    private readonly IProfileOutboxStorage _outboxStorage;
    private readonly IMessagingService _messagingService;

    public OutboxPublisher(
        ILogger<OutboxPublisher> logger,
        IProfileOutboxStorage outboxStorage,
        IMessagingService messagingService)
    {
        _logger = logger;
        _outboxStorage = outboxStorage;
        _messagingService = messagingService;
    }

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        _logger.LogInformation("Outbox processing started...");
        
        while (!stoppingToken.IsCancellationRequested)
        {
            try
            {
                _logger.LogDebug("Outbox processing iteration started...");
                
                var profiles = await _outboxStorage.GetProfileOutbox(10, stoppingToken);

                if (!profiles.Any())
                {
                    continue;
                }

                foreach (ProfileOutboxEntity profile in profiles)
                {
                    await _messagingService.Publish(new ProfileUpdatedMessage
                    {
                        ProfileId = profile.ProfileId,
                        Gender = profile.Gender,
                        Age = profile.Age,
                        Name = profile.Name,
                        Description = profile.Description,
                        Photos = profile.Photos
                    });
                }

                await _outboxStorage.ClearProfileOutbox(profiles.Select(x => x.OrderingId).ToList(), stoppingToken);
                
                _logger.LogDebug("Outbox part was cleared.");
            }
            catch (TaskCanceledException ex)
            {
                _logger.LogError(ex, "Operation was canceled.");
            }
            catch (Exception ex)
            {
                // todo: add circuit breaker for Postgresql or Kafka unavailability
                _logger.LogError(ex, "Error during outbox processing.");
            }
            finally
            {
                _logger.LogDebug("Outbox processing iteration was finished.");
                await Task.Delay(TimeSpan.FromSeconds(2), stoppingToken);
            }
        }
    }
}