using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using JwtGenerator;
using Microsoft.IdentityModel.Tokens;

var claims = new List<Claim> { new("ProfileId", 100.ToString()) };

var jwt = new JwtSecurityToken(
    issuer: AuthOptions.Issuer,
    audience: AuthOptions.Audience,
    claims: claims,
    expires: DateTime.UtcNow.Add(TimeSpan.FromDays(365)),
    signingCredentials: new SigningCredentials(AuthOptions.GetSymmetricSecurityKey(), SecurityAlgorithms.HmacSha256));

Console.WriteLine(new JwtSecurityTokenHandler().WriteToken(jwt));