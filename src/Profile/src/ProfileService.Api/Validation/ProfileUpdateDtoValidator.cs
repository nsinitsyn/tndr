using FluentValidation;
using TinderApiV1;

namespace ProfileService.Api.Validation;

public class ProfileUpdateDtoValidator : AbstractValidator<ProfileUpdateDto>
{
    public ProfileUpdateDtoValidator()
    {
        RuleFor(request => request.Age).GreaterThan(x => 17).LessThan(x => 100);
        RuleFor(request => request.Name).NotEmpty().MaximumLength(1000);
        RuleFor(request => request.Description).MaximumLength(10000);
    }
}