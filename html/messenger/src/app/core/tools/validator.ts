import {AbstractControl, ValidatorFn} from "@angular/forms";

export function customValidator<T>(validationFn: (value: T) => boolean): ValidatorFn {
  return (control: AbstractControl): { [key: string]: any } | null => {
    const value = control.value;

    // Use the provided validation function to check validity
    if (!validationFn(value)) {
      // Return error if the validation function returns false
      return {'customValidation': {valid: false}};
    }
    return null; // Return null if no error
  };
}
