using Grpc.Core;
using Grpc.Core.Interceptors;

namespace ProfileService.Api.Interceptors;

public class ExceptionHandlerInterceptor : Interceptor
{
    private readonly ILogger<ExceptionHandlerInterceptor> _logger;

    public ExceptionHandlerInterceptor(ILogger<ExceptionHandlerInterceptor> logger)
    {
        _logger = logger;
    }
    
    public override async Task<TResponse> UnaryServerHandler<TRequest, TResponse>(TRequest request, ServerCallContext context,
        UnaryServerMethod<TRequest, TResponse> continuation)
    {
        _logger.LogInformation(
            "Starting receiving call. Method: {method}",
            context.Method);
        
        try
        {
            return await base.UnaryServerHandler(request, context, continuation);
        }
        catch (Exception ex)
        {
            _logger.LogError(
                ex,
                "Error thrown by {method}. Headers: @{headers}. Request: @{request}",
                context.Method,
                context.RequestHeaders, request);
            
            throw new RpcException(new Status(StatusCode.Internal, ex.Message));
        }
    }
}