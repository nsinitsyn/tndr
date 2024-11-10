using System.Text.Json;
using Confluent.Kafka;
using ProfileService.Services.Entities.Messaging;

namespace ProfileService.Infrastructure.Messaging;

public class MessageSerializer : ISerializer<ProfileUpdatedMessage>
{
    public byte[] Serialize(ProfileUpdatedMessage data, SerializationContext context)
    {
        return JsonSerializer.SerializeToUtf8Bytes(data);
    }
}
