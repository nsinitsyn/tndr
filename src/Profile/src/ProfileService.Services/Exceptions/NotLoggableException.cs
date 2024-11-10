namespace ProfileService.Services.Exceptions;

/// <summary>
/// Exception that had already logged and should only handled for sending error to clients.
/// </summary>
public class NotLoggableException : Exception
{
    
}