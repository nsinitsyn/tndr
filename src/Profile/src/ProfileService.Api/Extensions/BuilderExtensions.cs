using Calzolari.Grpc.AspNetCore.Validation;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.IdentityModel.Tokens;
using ProfileService.Api.Authentication;
using ProfileService.Api.Health;
using ProfileService.Api.Interceptors;
using ProfileService.Api.Validation;
using ProfileService.Infrastructure.Configuration;
using ProfileService.Infrastructure.Extensions;
using ProfileService.Services.Extensions;
using Serilog;

namespace ProfileService.Api.Extensions;

public static class BuilderExtensions
{
    public static void ConfigureBuilder(this WebApplicationBuilder builder)
    {
        builder.Host.UseSerilog(
            (context, configuration) => configuration.ReadFrom.Configuration(context.Configuration));
        
        var kafkaSection = builder.Configuration.GetSection(nameof(KafkaConfiguration));
        if (!kafkaSection.Exists())
        {
            throw new InvalidOperationException(
                $"The section {nameof(KafkaConfiguration)} wasn't found in configuration.");
        }
        builder.Services.Configure<KafkaConfiguration>(settings => kafkaSection.Bind(settings));
        
        builder.Services.AddEndpointsApiExplorer();
        builder.Services.AddSwaggerGen();
        
        builder.Services.AddHttpContextAccessor();
        
        builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
            .AddJwtBearer(options =>
            {
                options.TokenValidationParameters = new TokenValidationParameters
                {
                    ValidateIssuer = true,
                    ValidIssuer = AuthOptions.Issuer,
                    ValidateAudience = true,
                    ValidAudience = AuthOptions.Audience,
                    ValidateLifetime = true,
                    IssuerSigningKey = AuthOptions.GetSymmetricSecurityKey(),
                    ValidateIssuerSigningKey = true,
                };
            });
        builder.Services.AddAuthorization();

        builder.Services.AddGrpc(options =>
        {
            options.EnableMessageValidation();
            options.Interceptors.Add<ExceptionHandlerInterceptor>();
            options.Interceptors.Add<UserProfileInterceptor>();
        });
        builder.Services.AddGrpcReflection();
        builder.Services.AddGrpcValidation();
        builder.Services.AddValidator<UpdateProfileRequestValidator>();
        builder.Services.AddValidator<ProfileUpdateDtoValidator>();

        var kafkaConfiguration = kafkaSection.Get<KafkaConfiguration>();
        ArgumentNullException.ThrowIfNull(kafkaConfiguration);

        builder.Services.AddHealthChecks()
            .AddCheck<ReadinessHealthCheck>("Startup", tags: new[] { "ready" })
            // todo: need it?
            .AddKafka(options =>
            {
                options.BootstrapServers = kafkaConfiguration.BootstrapServers;
            });

        var connectionString = builder.Configuration.GetConnectionString("Profile");
        if (string.IsNullOrWhiteSpace(connectionString))
        {
            throw new InvalidOperationException("Connection string cannot be null or empty.");
        }
        builder.Services.AddNpgsqlDataSource(connectionString);
        
        builder.Services.AddScoped<UserProfileProvider>();
        builder.Services.AddScoped<IUserProfileSetter>(sp => sp.GetRequiredService<UserProfileProvider>());
        builder.Services.AddScoped<IUserProfileProvider>(sp => sp.GetRequiredService<UserProfileProvider>());

        builder.Services.AddApplicationServices();
        builder.Services.AddInfrastructureServices();
    }
}
