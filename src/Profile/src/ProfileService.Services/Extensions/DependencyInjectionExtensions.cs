using Microsoft.Extensions.DependencyInjection;

namespace ProfileService.Services.Extensions;

public static class DependencyInjectionExtensions
{
    public static void AddApplicationServices(this IServiceCollection services)
    {
        services.AddScoped<ProfileService>();
        services.AddHostedService<OutboxPublisher>();
    }
}
