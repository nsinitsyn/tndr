using FluentMigrator.Runner.VersionTableInfo;

namespace ProfileService.DbMigrator;

[VersionTableMetaData]
public class CustomVersionTableMetaData : IVersionTableMetaData
{
    public virtual string SchemaName => "public";

    public virtual string TableName => "version_info";

    public virtual string ColumnName => "version";

    public virtual string UniqueIndexName => "version_info_version_key";

    public virtual string AppliedOnColumnName => "applied_on";
    
    public virtual string DescriptionColumnName => "description";

    public virtual bool OwnsSchema => false;

    public bool CreateWithPrimaryKey { get; } = false;
}