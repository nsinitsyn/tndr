using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Microsoft.Extensions.Diagnostics.HealthChecks;
using ProfileService.Api.GrpcServices;
using Serilog;

namespace ProfileService.Api.Extensions;

public static class WebApplicationExtensions
{
    public static void ConfigureApplication(this WebApplication app)
    {
        TaskScheduler.UnobservedTaskException += (_, e) =>
            app.Services.GetRequiredService<ILogger<Program>>().LogError(e.Exception, "Unhandled exception.");

        app.UseAuthentication();
        app.UseAuthorization();
        
        app.UseSerilogRequestLogging();
        
        if (app.Environment.IsDevelopment())
        {
            app.UseSwagger();
            app.UseSwaggerUI();
            app.MapGrpcReflectionService();
        }

        app.MapGrpcService<ProfileGrpcService>();

        app.MapHealthChecks("/healthz/ready", new HealthCheckOptions
        {
            Predicate = healthCheck => healthCheck.Tags.Contains("ready"),
        });

        app.MapHealthChecks("/healthz/live", new HealthCheckOptions
        {
            Predicate = _ => false,
            ResultStatusCodes =
            {
                [HealthStatus.Healthy] = StatusCodes.Status200OK,
                [HealthStatus.Degraded] = StatusCodes.Status200OK,
                [HealthStatus.Unhealthy] = StatusCodes.Status503ServiceUnavailable
            }
        });
    }
}
