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
        
        builder.Services.AddHealthChecks()
            .AddCheck<ReadinessHealthCheck>("Startup", tags: new[] { "ready" });
        
        AddAuthentication(builder.Services);
        AddAuthorization(builder.Services);
        AddGrpc(builder.Services);
        AddNpgsql(builder.Services, builder.Configuration);

        AddServices(builder.Services);
        builder.Services.AddApplicationServices();
        builder.Services.AddInfrastructureServices();
    }
    
    private static void AddServices(IServiceCollection services)
    {
        services.AddScoped<UserProfileProvider>();
        services.AddScoped<IUserProfileSetter>(sp => sp.GetRequiredService<UserProfileProvider>());
        services.AddScoped<IUserProfileProvider>(sp => sp.GetRequiredService<UserProfileProvider>());
    }

    private static void AddAuthentication(IServiceCollection services)
    {
        services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
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
    }

    private static void AddAuthorization(IServiceCollection services)
    {
        services.AddAuthorization(options =>
        {
            options.AddPolicy("Administrator", policy =>
            {
                policy.AuthenticationSchemes.Add(JwtBearerDefaults.AuthenticationScheme);
                policy.RequireClaim("Role", "Administrator");
            });
            
            options.AddPolicy("MatchService", policy =>
            {
                policy.AuthenticationSchemes.Add(JwtBearerDefaults.AuthenticationScheme);
                policy.RequireClaim("Service", "MatchService");
            });

            options.AddPolicy("AdministratorOrMatchService", policy =>
            {
                policy.AuthenticationSchemes.Add(JwtBearerDefaults.AuthenticationScheme);
                policy.RequireAssertion(context =>
                    context.User.HasClaim(c =>
                        c is { Type: "Role", Value: "Administrator" } or { Type: "Service", Value: "MatchService" }));
            });
        });
    }
    
    private static void AddGrpc(IServiceCollection services)
    {
        services.AddGrpc(options =>
        {
            options.EnableMessageValidation();
            options.Interceptors.Add<UserProfileInterceptor>();
            options.Interceptors.Add<LoggingInterceptor>();
        });
        services.AddGrpcReflection();
        services.AddGrpcValidation();
        services.AddValidator<ProfileDtoValidator>();
        services.AddValidator<CreateProfileValidator>();
        services.AddValidator<UpdateProfileValidator>();
    }

    private static void AddNpgsql(IServiceCollection services, IConfiguration configuration)
    {
        var connectionString = configuration.GetConnectionString("Profile");
        if (string.IsNullOrWhiteSpace(connectionString))
        {
            throw new InvalidOperationException("Connection string cannot be null or empty.");
        }
        services.AddNpgsqlDataSource(connectionString);
    }
}
