using FluentValidation;
using TinderApiV1;

namespace ProfileService.Api.Validation;

public class CreateProfileValidator : AbstractValidator<CreateProfileRequest>
{
    public CreateProfileValidator()
    {
        RuleFor(request => request.Profile).NotNull();
        RuleFor(customer => customer.Profile).SetValidator(new ProfileDtoValidator());
    }
}
