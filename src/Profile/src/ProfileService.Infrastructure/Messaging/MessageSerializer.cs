using System.Text.Json;
using System.Text.Json.Serialization;
using Confluent.Kafka;
using ProfileService.Services.Entities.Messaging;

namespace ProfileService.Infrastructure.Messaging;

public class MessageSerializer : ISerializer<ProfileUpdatedMessage>
{
    public byte[] Serialize(ProfileUpdatedMessage data, SerializationContext context)
    {
        var options = new JsonSerializerOptions();
        options.Converters.Add(new JsonStringEnumConverter());
        
        return JsonSerializer.SerializeToUtf8Bytes(data, options);
    }
}
