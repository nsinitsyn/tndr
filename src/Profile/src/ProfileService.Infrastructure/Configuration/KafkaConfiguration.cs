namespace ProfileService.Infrastructure.Configuration;

public class KafkaConfiguration
{
    public string BootstrapServers { get; set; } = string.Empty;

    public string TopicName { get; set; } = string.Empty;
}