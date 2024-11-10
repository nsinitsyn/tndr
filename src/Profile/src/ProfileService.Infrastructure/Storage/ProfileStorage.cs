using Npgsql;
using NpgsqlTypes;
using ProfileService.Domain;
using ProfileService.Services.Dependencies;

namespace ProfileService.Infrastructure.Storage;

public class ProfileStorage : IProfileStorage
{
    private readonly NpgsqlDataSource _dataSource;

    public ProfileStorage(NpgsqlDataSource dataSource)
    {
        _dataSource = dataSource;
    }
    
    public async Task AddProfile(ProfileEntity profile)
    {
        await using var connection = _dataSource.CreateConnection();
        await connection.OpenAsync();
        
        await using var command = new NpgsqlCommand(
            """
            WITH insert_profile AS (
                INSERT INTO profile (sex, age, name, description, photos)
                VALUES ($1, $2, $3, $4, $5)
                RETURNING id, sex, age, name, description, photos
            )
            INSERT INTO profile_outbox(profile_id, sex, age, name, description, photos)
            SELECT * FROM insert_profile
            """, connection)
        {
            Parameters =
            {
                new() { Value = profile.Sex, NpgsqlDbType = NpgsqlDbType.Boolean },
                new() { Value = profile.Age, NpgsqlDbType = NpgsqlDbType.Boolean },
                new() { Value = profile.Name, NpgsqlDbType = NpgsqlDbType.Boolean },
                new() { Value = profile.Description, NpgsqlDbType = NpgsqlDbType.Boolean },
            }
        };
        
        await command.ExecuteNonQueryAsync();
    }

    public async Task UpdateProfile(ProfileEntity profile)
    {
        await using var connection = _dataSource.CreateConnection();
        await connection.OpenAsync();
        await using var transaction = await connection.BeginTransactionAsync();
        
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

        await insertProfileCommand.ExecuteNonQueryAsync();
        
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

        await insertOutboxCommand.ExecuteNonQueryAsync();

        await transaction.CommitAsync();
    }

    public Task<ProfileEntity> GetProfiles()
    {
        throw new NotImplementedException();
    }
}