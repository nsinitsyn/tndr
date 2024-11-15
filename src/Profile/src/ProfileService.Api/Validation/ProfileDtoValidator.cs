using FluentValidation;
using TinderApiV1;

namespace ProfileService.Api.Validation;

public class ProfileDtoValidator : AbstractValidator<ProfileDto>
{
    public ProfileDtoValidator()
    {
        RuleFor(request => request.Age).GreaterThan(x => 17).LessThan(x => 100);
        RuleFor(request => request.Name).NotEmpty().MaximumLength(1000);
        RuleFor(request => request.Description).MaximumLength(10000);
        RuleFor(request => request.Gender)
            .Must(x => !x.Equals(Gender.Unspecified))
            .WithMessage("Gender should be specified.");
    }
}