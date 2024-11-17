using Confluent.Kafka;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using ProfileService.Infrastructure.Configuration;
using ProfileService.Services.Dependencies;
using ProfileService.Services.Entities.Messaging;

namespace ProfileService.Infrastructure.Messaging;

public class MessagingService : IMessagingService, IDisposable
{
    private readonly ILogger<MessagingService> _logger;
    private readonly KafkaConfiguration _configuration;
    private readonly IProducer<string, ProfileUpdatedMessage> _producer;

    public MessagingService(
        ILogger<MessagingService> logger,
        IOptions<KafkaConfiguration> configuration)
    {
        _logger = logger;
        _configuration = configuration.Value;
        
        var producerConfig = new ProducerConfig
        {
            BootstrapServers = _configuration.BootstrapServers,
        };

        _producer = new ProducerBuilder<string, ProfileUpdatedMessage>(producerConfig)
            .SetValueSerializer(new MessageSerializer())
            .Build();
    }
    
    public async Task Publish(ProfileUpdatedMessage message)
    {
        await _producer.ProduceAsync(
            _configuration.TopicName,
            new Message<string, ProfileUpdatedMessage> { Key = message.Id.ToString(), Value = message, });
        
        _logger.LogInformation("Message was published: {@message}", message);
    }

    public void Dispose()
    {
        _producer.Dispose();
    }
}