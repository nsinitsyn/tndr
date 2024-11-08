using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.IdentityModel.Tokens;
using ProfileService.Api.Authentication;
using ProfileService.Api.GrpcInterceptors;
using ProfileService.Api.Health;
using ProfileService.Services.Extensions;
using Serilog;

namespace ProfileService.Api.Extensions;

public static class WebApplicationBuilderExtensions
{
    public static void ConfigureBuilder(this WebApplicationBuilder builder)
    {
        builder.Host.UseSerilog(
            (context, configuration) => configuration.ReadFrom.Configuration(context.Configuration));
        
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
            options.Interceptors.Add<ExceptionHandlerInterceptor>();
            options.Interceptors.Add<UserProfileInterceptor>();
        });
        builder.Services.AddGrpcReflection();
        
        builder.Services.AddHealthChecks().AddCheck<ReadinessHealthCheck>("Startup", tags: new[] { "ready" });

        builder.Services.AddScoped<UserProfileProvider>();
        builder.Services.AddScoped<IUserProfileSetter>(sp => sp.GetRequiredService<UserProfileProvider>());
        builder.Services.AddScoped<IUserProfileProvider>(sp => sp.GetRequiredService<UserProfileProvider>());

        builder.Services.AddApplicationServices();
    }
}
