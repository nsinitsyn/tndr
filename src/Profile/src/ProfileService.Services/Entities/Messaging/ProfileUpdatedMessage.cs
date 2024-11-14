namespace ProfileService.Services.Entities.Messaging;

public class ProfileUpdatedMessage
{
    public long ProfileId { get; set; }
    public char Gender { get; set; }
    public int Age { get; set; }
    public string Name { get; set; } = null!;
    public string Description { get; set; } = null!;
    public string[] Photos { get; set; } = [];
}