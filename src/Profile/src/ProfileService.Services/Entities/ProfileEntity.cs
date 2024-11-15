namespace ProfileService.Services.Entities;

public class ProfileEntity
{
    public long ProfileId { get; set; }
    public Gender Gender { get; set; }
    public int Age { get; set; }
    public string Name { get; set; } = null!;
    public string Description { get; set; } = null!;
    public string[] Photos { get; set; } = [];
}