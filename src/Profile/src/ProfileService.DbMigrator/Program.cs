using FluentMigrator.Runner;
using FluentMigrator.Runner.VersionTableInfo;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using ProfileService.DbMigrator;
using ProfileService.DbMigrator.Migrations;

var builder = new ConfigurationBuilder();

builder.SetBasePath(Directory.GetCurrentDirectory())
    .AddJsonFile("appsettings.json", optional: false, reloadOnChange: true);

IConfiguration config = builder.Build();
var pConnectionString = config.GetConnectionString("Profile");

if (string.IsNullOrWhiteSpace(pConnectionString))
{
    Console.WriteLine("Connection string cannot be empty.");
    return;
}

using (var serviceProvider = CreateServices(pConnectionString))
{
    using (var scope = serviceProvider.CreateScope())
    {
        UpdateDatabase(scope.ServiceProvider);
    }
}

Console.WriteLine("Finished. Press any key for exit...");
Console.ReadKey();

ServiceProvider CreateServices(string connectionString)
{
    return new ServiceCollection()
        .AddFluentMigratorCore()
        .ConfigureRunner(rb => rb
            .AddPostgres()
            .WithGlobalConnectionString(connectionString)
            .ScanIn(typeof(M000_InitialCreate).Assembly).For.Migrations())
        .AddLogging(lb => lb.AddFluentMigratorConsole())
        .AddScoped(typeof(IVersionTableMetaData), typeof(CustomVersionTableMetaData))
        .BuildServiceProvider(false);
}

void UpdateDatabase(IServiceProvider serviceProvider)
{
    var runner = serviceProvider.GetRequiredService<IMigrationRunner>();
    runner.MigrateUp();
}