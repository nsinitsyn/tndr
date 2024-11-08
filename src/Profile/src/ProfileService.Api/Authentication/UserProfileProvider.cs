namespace ProfileService.Api.Authentication;

public interface IUserProfileProvider
{
    long ProfileId { get; }
}

public interface IUserProfileSetter
{
    long ProfileId { set; }
}

public class UserProfileProvider : IUserProfileProvider, IUserProfileSetter 
{
    public long ProfileId { get; set; }
}