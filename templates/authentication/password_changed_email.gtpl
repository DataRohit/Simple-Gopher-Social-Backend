{{- /* gotype: gopher-social-backend-server/pkg/mailer/passwordChangedMail.PasswordChangedEmailData */ -}}

<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Successfully Reset</title>
</head>
<body>
    <h1>Password Reset Successful!</h1>
    <p>Hi {{.Email}},</p>
    <p>Your password has been successfully reset. You can now log in using your new password.</p>
    <p>Thank you!</p>
</body>
</html>
