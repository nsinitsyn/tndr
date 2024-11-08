using System.Text;
using Microsoft.IdentityModel.Tokens;

namespace JwtGenerator;

public class AuthOptions
{
    public const string Issuer = "AuthServer";
    public const string Audience = "AuthClient";
    private const string Key = "fjg847sdjvnjxcFHdsag38d_d8sj3aqQwfdsph3456v0bjz45ty54gpo3vhjs7234f09Odp";
    public static SymmetricSecurityKey GetSymmetricSecurityKey() => 
        new SymmetricSecurityKey(Encoding.UTF8.GetBytes(Key));
}