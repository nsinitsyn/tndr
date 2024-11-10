using Microsoft.Extensions.Logging;
using Npgsql;
using NpgsqlTypes;
using ProfileService.Services.Dependencies;
using ProfileService.Services.Entities;
using ProfileService.Services.Exceptions;

namespace ProfileService.Infrastructure.Storage;

public class ProfileStorage : IProfileStorage
{
    private readonly ILogger<ProfileStorage> _logger;
    private readonly NpgsqlDataSource _dataSource;

    public ProfileStorage(ILogger<ProfileStorage> logger, NpgsqlDataSource dataSource)
    {
        _logger = logger;
        _dataSource = dataSource;
    }

    public async Task<ProfileEntity> GetProfile(long profileId, CancellationToken cancellationToken)
    {
        await using var connection = _dataSource.CreateConnection();
        await connection.OpenAsync(cancellationToken);

        await using var command = new NpgsqlCommand(
            """
            SELECT sex, age, name, description, photos
            FROM profile
            WHERE id = $1
            """, connection);
        
        command.Parameters.Add(new() { Value = profileId });
        
        var reader = await command.ExecuteReaderAsync(cancellationToken);

        if (!await reader.ReadAsync(cancellationToken))
        {
            _logger.LogError("Unexpected behavior. Get profile result is null. Profile id: {profileId}", profileId);
            throw new NotLoggableException();
        }

        return new ProfileEntity
        {
            ProfileId = profileId,
            Sex = reader.GetBoolean(0),
            Age = reader.GetInt16(1),
            Name = reader.GetString(2),
            Description = reader.GetString(3),
            Photos = (string[])reader.GetValue(4)
        };
    }
    
    public async Task<long> CreateProfile(CreateProfileEntity profile, CancellationToken cancellationToken)
    {
        await using var connection = _dataSource.CreateConnection();
        await connection.OpenAsync(cancellationToken);

        await using var command = new NpgsqlCommand(
            """
            WITH insert_profile AS (
                INSERT INTO profile (sex, age, name, description, photos)
                VALUES ($1, $2, $3, $4, $5)
                RETURNING id, sex, age, name, description, photos
            )
            INSERT INTO profile_outbox(profile_id, sex, age, name, description, photos)
            SELECT * FROM insert_profile
            RETURNING profile_id
            """, connection);
        
        command.Parameters.Add(new() { Value = profile.Sex });
        command.Parameters.Add(new() { Value = profile.Age });
        command.Parameters.Add(new() { Value = profile.Name });
        command.Parameters.Add(new() { Value = profile.Description });
        command.Parameters.Add(new() { Value = profile.Photos, NpgsqlDbType = NpgsqlDbType.Array | NpgsqlDbType.Text });
        
        var executionResult = await command.ExecuteScalarAsync(cancellationToken);
        if (executionResult == null)
        {
            _logger.LogError("Unexpected behavior. Insert profile result is null. Profile: {@profile}", profile);
            throw new NotLoggableException();
        }
        
        var profileId = long.Parse(executionResult.ToString()!);
        return profileId;
    }

    public async Task UpdateProfile(ProfileEntity profile, CancellationToken cancellationToken)
    {
        await using var connection = _dataSource.CreateConnection();
        await connection.OpenAsync(cancellationToken);
        await using var transaction = await connection.BeginTransactionAsync(cancellationToken);
        
        await using var insertProfileCommand = new NpgsqlCommand(
            """
            UPDATE profile 
            SET
            	sex = $2,
            	age = $3,
            	name = $4,
            	description = $5,
            	photos = $6,
            	changed_at = timezone('utc', now())
            WHERE id = $1;
            """, connection, transaction);
        
        insertProfileCommand.Parameters.Add(new() { Value = profile.ProfileId });
        insertProfileCommand.Parameters.Add(new() { Value = profile.Sex });
        insertProfileCommand.Parameters.Add(new() { Value = profile.Age });
        insertProfileCommand.Parameters.Add(new() { Value = profile.Name });
        insertProfileCommand.Parameters.Add(new() { Value = profile.Description });
        insertProfileCommand.Parameters.Add(new() { Value = profile.Photos, NpgsqlDbType = NpgsqlDbType.Array | NpgsqlDbType.Text });

        await insertProfileCommand.ExecuteNonQueryAsync(cancellationToken);
        
        await using var insertOutboxCommand = new NpgsqlCommand(
            """
            INSERT INTO profile_outbox (profile_id, sex, age, name, description, photos)
            VALUES ($1, $2, $3, $4, $5, $6)
            """, connection, transaction);
        
        insertOutboxCommand.Parameters.Add(new() { Value = profile.ProfileId });
        insertOutboxCommand.Parameters.Add(new() { Value = profile.Sex });
        insertOutboxCommand.Parameters.Add(new() { Value = profile.Age });
        insertOutboxCommand.Parameters.Add(new() { Value = profile.Name });
        insertOutboxCommand.Parameters.Add(new() { Value = profile.Description });
        insertOutboxCommand.Parameters.Add(new() { Value = profile.Photos, NpgsqlDbType = NpgsqlDbType.Array | NpgsqlDbType.Text });

        await insertOutboxCommand.ExecuteNonQueryAsync(cancellationToken);

        await transaction.CommitAsync(cancellationToken);
    }

    public async Task<IList<ProfileEntity>> GetProfiles(IList<long> profileIds, CancellationToken cancellationToken)
    {
        var result = new List<ProfileEntity>();
        
        await using var connection = _dataSource.CreateConnection();
        await connection.OpenAsync(cancellationToken);
        
        var command = new NpgsqlCommand(
            """
            SELECT id, sex, age, name, description, photos 
            FROM profile
            WHERE id = ANY(:ids)
            """,
            connection
        );

        command.Parameters.Add(new()
            { ParameterName = "ids", Value = profileIds, NpgsqlDbType = NpgsqlDbType.Array | NpgsqlDbType.Bigint });
        
        var reader = await command.ExecuteReaderAsync(cancellationToken);
        while (await reader.ReadAsync(cancellationToken))
        {
            result.Add(new ProfileEntity
            {
                ProfileId = reader.GetInt64(0),
                Sex = reader.GetBoolean(1),
                Age = reader.GetInt16(2),
                Name = reader.GetString(3),
                Description = reader.GetString(4),
                Photos = (string[])reader.GetValue(5)
            });
        }

        return result;
    }
}