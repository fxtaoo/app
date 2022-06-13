Import-Module ActiveDirectory

$maxPwdAge=(Get-ADDefaultDomainPasswordPolicy).MaxPasswordAge.Days
$startMonitor=(get-date).AddDays(14-$maxPwdAge).ToShortDateString()


Get-ADUser -filter {Enabled -eq $True -and PasswordNeverExpires -eq $False -and PasswordLastSet -gt 0} –Properties * | where {($_.PasswordLastSet).ToShortDateString() -lt $startMonitor} | select-object SamAccountName,passwordlastset,Name | Export-Csv -Path ./userdata.csv -encoding utf8 -NoTypeInformation