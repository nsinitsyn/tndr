using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using JwtGenerator;
using Microsoft.IdentityModel.Tokens;

var claims = new List<Claim> { new("ProfileId", 1.ToString()), new("Gender", "M") };
// var claims = new List<Claim> { new("Role", "User") };

var jwt = new JwtSecurityToken(
    issuer: AuthOptions.Issuer,
    audience: AuthOptions.Audience,
    claims: claims,
    expires: DateTime.UtcNow.Add(TimeSpan.FromDays(365)),
    signingCredentials: new SigningCredentials(AuthOptions.GetSymmetricSecurityKey(), SecurityAlgorithms.HmacSha256));

Console.WriteLine(new JwtSecurityTokenHandler().WriteToken(jwt));

// ProfileId: 1, Gender: M
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQcm9maWxlSWQiOiIxIiwiR2VuZGVyIjoiTSIsImV4cCI6MTc2MzIwNzQxMywiaXNzIjoiQXV0aFNlcnZlciIsImF1ZCI6IkF1dGhDbGllbnQifQ.VAVP65lIUhabxR4UknvQkRKiVCfu116cf3tZC8-dsfw

// ProfileId: 10
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQcm9maWxlSWQiOiIxMCIsImV4cCI6MTc2MjcxOTY2NSwiaXNzIjoiQXV0aFNlcnZlciIsImF1ZCI6IkF1dGhDbGllbnQifQ.pDWIoPzTE9Q_ccKgC11CMiczkKx52dYikYZEC6qvAbU

// Role: Administrator
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJSb2xlIjoiQWRtaW5pc3RyYXRvciIsImV4cCI6MTc2Mjc4MjM5OCwiaXNzIjoiQXV0aFNlcnZlciIsImF1ZCI6IkF1dGhDbGllbnQifQ.rVqdEV044q5b4QtXu4fJsGsyh1bFk8zp_EiPFwzh5LE

// Role: User
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJSb2xlIjoiVXNlciIsImV4cCI6MTc2Mjc4Mjg1NCwiaXNzIjoiQXV0aFNlcnZlciIsImF1ZCI6IkF1dGhDbGllbnQifQ.CLHHAutq5IM4zMNZHlvSxJU9H9C3N9ruovfMdGfpmFo