using FluentMigrator;

namespace ProfileService.DbMigrator.Migrations;

public class M000_InitialCreate
{
    [Migration(0)]
    public class InitialCreate : Migration
    {
        public override void Up()
        {
            CreateProfileTable();
            CreateProfileOutboxTable();
        }

        public override void Down()
        {
            Delete.Table("profile");
            Delete.Table("profile_outbox");
        }
        
        private void CreateProfileTable()
        {
            Create.Table("profile")
                .WithColumn("id").AsInt64().PrimaryKey().Identity()
                .WithColumn("sex").AsBoolean().NotNullable()
                .WithColumn("age").AsInt16().NotNullable()
                .WithColumn("name").AsString().NotNullable()
                .WithColumn("description").AsString().NotNullable()
                .WithColumn("photos").AsCustom("text[]").NotNullable()
                .WithColumn("created_at").AsDateTime2().NotNullable().WithDefault(SystemMethods.CurrentUTCDateTime)
                .WithColumn("changed_at").AsDateTime2().NotNullable().WithDefault(SystemMethods.CurrentUTCDateTime);
        }
        
        private void CreateProfileOutboxTable()
        {
            Create.Table("profile_outbox")
                .WithColumn("ordering_id").AsInt64().PrimaryKey().Identity()
                .WithColumn("profile_id").AsInt64()
                .WithColumn("sex").AsBoolean().NotNullable()
                .WithColumn("age").AsInt16().NotNullable()
                .WithColumn("name").AsString().NotNullable()
                .WithColumn("description").AsString().NotNullable()
                .WithColumn("photos").AsCustom("text[]").NotNullable();
        }
    }
}