using Grpc.Core;
using Grpc.Core.Interceptors;
using ProfileService.Api.Authentication;

namespace ProfileService.Api.Interceptors;

public class UserProfileInterceptor : Interceptor
{
    private readonly IHttpContextAccessor _httpContextAccessor;
    private readonly IUserProfileSetter _userProfileSetter;
    private readonly ILogger<UserProfileInterceptor> _logger;

    public UserProfileInterceptor(
        IHttpContextAccessor httpContextAccessor,
        IUserProfileSetter userProfileSetter,
        ILogger<UserProfileInterceptor> logger)
    {
        _httpContextAccessor = httpContextAccessor;
        _userProfileSetter = userProfileSetter;
        _logger = logger;
    }

    public override Task<TResponse> UnaryServerHandler<TRequest, TResponse>(TRequest request, ServerCallContext context,
        UnaryServerMethod<TRequest, TResponse> continuation)
    {
        if (_httpContextAccessor.HttpContext == null)
        {
            return base.UnaryServerHandler(request, context, continuation);
        }
        
        var userPrincipal = _httpContextAccessor.HttpContext.User;
        
        var printableClaims = string.Join(", ", userPrincipal.Claims.Select(x => $"{x.Type}:{x.Value}"));
        _logger.LogInformation("Request claims: {claims}", printableClaims);
        
        var profileIdClaim = userPrincipal.FindFirst(CustomClaimTypes.ProfileId);
        
        if (profileIdClaim != null)
        {
            if (long.TryParse(profileIdClaim.Value, out var profileId))
            {
                _userProfileSetter.ProfileId = profileId;
            }
            else
            {
                _logger.LogWarning(
                    "Cannot parse profile id claim from jwt token: {profileIdClaim}.",
                    profileIdClaim.Value);
            }
        }
        
        return base.UnaryServerHandler(request, context, continuation);
    }
}