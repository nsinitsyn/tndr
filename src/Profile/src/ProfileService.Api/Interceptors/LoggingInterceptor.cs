using Grpc.Core;
using Grpc.Core.Interceptors;
using ProfileService.Api.Authentication;
using ProfileService.Services.Exceptions;

namespace ProfileService.Api.Interceptors;

public class LoggingInterceptor : Interceptor
{
    private readonly ILogger<LoggingInterceptor> _logger;
    private readonly IUserProfileProvider _userProfileProvider;

    public LoggingInterceptor(ILogger<LoggingInterceptor> logger, IUserProfileProvider userProfileProvider)
    {
        _logger = logger;
        _userProfileProvider = userProfileProvider;
    }
    
    public override async Task<TResponse> UnaryServerHandler<TRequest, TResponse>(TRequest request, ServerCallContext context,
        UnaryServerMethod<TRequest, TResponse> continuation)
    {
        _logger.LogInformation(
            "Starting receiving call. Method: {method}. ProfileId: {profileId}. Headers: {@headers}. Request: {@request}.",
            context.Method, _userProfileProvider.ProfileId, context.RequestHeaders, request);

        try
        {
            var result = await base.UnaryServerHandler(request, context, continuation);
            
            _logger.LogInformation(
                "Finished call. Result: {@result}.",
                result);

            return result;
        }
        catch (NotLoggableException)
        {
            throw new RpcException(new Status(StatusCode.Internal, "Internal server error."));
        }
        catch (TaskCanceledException ex)
        {
            _logger.LogWarning(
                ex,
                "Task cancelled in {method}.",
                context.Method);

            throw new RpcException(new Status(StatusCode.Cancelled, string.Empty));
        }
        catch (Exception ex)
        {
            _logger.LogError(
                ex,
                "Error thrown by {method}.",
                context.Method);
            
            throw new RpcException(new Status(StatusCode.Internal, ex.Message));
        }
    }
}