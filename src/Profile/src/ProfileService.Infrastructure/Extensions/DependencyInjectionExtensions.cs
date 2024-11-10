using Microsoft.Extensions.DependencyInjection;
using ProfileService.Infrastructure.Messaging;
using ProfileService.Infrastructure.Storage;
using ProfileService.Services.Dependencies;

namespace ProfileService.Infrastructure.Extensions;

public static class DependencyInjectionExtensions
{
    public static void AddInfrastructureServices(this IServiceCollection services)
    {
        services.AddTransient<IProfileStorage, ProfileStorage>();
        services.AddTransient<IProfileOutboxStorage, ProfileOutboxStorage>();
        services.AddSingleton<IMessagingService, MessagingService>();
    }
}
