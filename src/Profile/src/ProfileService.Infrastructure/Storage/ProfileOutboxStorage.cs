using Npgsql;
using NpgsqlTypes;
using ProfileService.Services.Dependencies;
using ProfileService.Services.Entities;

namespace ProfileService.Infrastructure.Storage;

public class ProfileOutboxStorage : IProfileOutboxStorage
{
    private readonly NpgsqlDataSource _dataSource;

    public ProfileOutboxStorage(NpgsqlDataSource dataSource)
    {
        _dataSource = dataSource;
    }
    
    public async Task<IList<ProfileOutboxEntity>> GetProfileOutbox(int limit, CancellationToken cancellationToken)
    {
        var result = new List<ProfileOutboxEntity>(limit);
        
        await using var connection = _dataSource.CreateConnection();
        await connection.OpenAsync(cancellationToken);

        var command = new NpgsqlCommand(
            """
            SELECT ordering_id, profile_id, gender, age, name, description, photos
            FROM profile_outbox
            ORDER BY ordering_id
            LIMIT $1
            """,
            connection
        );
        
        command.Parameters.Add(new() { Value = limit });

        var reader = await command.ExecuteReaderAsync(cancellationToken);
        while (await reader.ReadAsync(cancellationToken))
        {
            result.Add(new ProfileOutboxEntity
            {
                OrderingId = reader.GetInt64(0),
                ProfileId = reader.GetInt64(1),
                Gender = reader.GetFieldValue<Gender>(2),
                Age = reader.GetInt16(3),
                Name = reader.GetString(4),
                Description = reader.GetString(5),
                Photos = (string[])reader.GetValue(6)
            });
        }

        return result;
    }

    public async Task ClearProfileOutbox(List<long> ids, CancellationToken cancellationToken)
    {
        await using var connection = _dataSource.CreateConnection();
        await connection.OpenAsync(cancellationToken);

        var command = new NpgsqlCommand(
            """
            DELETE FROM profile_outbox
            WHERE ordering_id = ANY(:ids)
            """,
            connection
        );

        command.Parameters.Add(new()
            { ParameterName = "ids", Value = ids, NpgsqlDbType = NpgsqlDbType.Array | NpgsqlDbType.Bigint });

        await command.ExecuteNonQueryAsync(cancellationToken);
    }
}