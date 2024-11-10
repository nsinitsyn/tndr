using FluentValidation;
using TinderApiV1;

namespace ProfileService.Api.Validation;

public class UpdateProfileRequestValidator : AbstractValidator<UpdateProfileRequest>
{
    public UpdateProfileRequestValidator()
    {
        RuleFor(request => request.Profile).NotNull();
        RuleFor(customer => customer.Profile).SetValidator(new ProfileUpdateDtoValidator());
    }
}
