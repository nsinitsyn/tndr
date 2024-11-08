using Grpc.Core;
using Grpc.Core.Interceptors;
using TinderApiV1;

namespace ProfileService.Api.GrpcInterceptors;

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
            _logger.LogError(ex, "Error thrown by {method}.", context.Method);
            
            // if TResponse contains UnexpectedServerError field then return TResponse instance.
            var responseType = continuation.Method.ReturnType.GenericTypeArguments.FirstOrDefault();
            if(responseType == null)
            {
                throw;
            }

            var property = responseType.GetProperty(nameof(UpdateProfileResponse.ErrorOneofCase.UnexpectedServerError));
            var response = Activator.CreateInstance(responseType);
            
            if (property != null && response != null)
            {
                property.SetValue(response, new UnexpectedServerError
                {
                    InternalDetails = ex.ToString(),
                    ExternalDetails = "Internal Server Error"
                });

                return (TResponse)response;
            }

            throw;
        }
    }
}