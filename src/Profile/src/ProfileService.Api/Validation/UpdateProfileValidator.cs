using FluentValidation;
using TinderApiV1;

namespace ProfileService.Api.Validation;

public class UpdateProfileValidator : AbstractValidator<UpdateProfileRequest>
{
    public UpdateProfileValidator()
    {
        RuleFor(request => request.Profile).NotNull();
        RuleFor(customer => customer.Profile).SetValidator(new ProfileDtoValidator());
    }
}
