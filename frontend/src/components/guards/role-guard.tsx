import { useAuthentication } from '../providers/authentication-provider';

export default function RoleGuard({
  value,
  children,
}: {
  value: string | string[];
  children: React.ReactNode;
}) {
  const { user, isLoading } = useAuthentication();

  if (isLoading) return null;
  if (!user) return null;

  if (Array.isArray(value)) {
    if (
      !value.includes(
        Array.isArray(user.Role) ? user.Role.join(', ') : (user.Role ?? 'user')
      )
    )
      return null;
  } else {
    if (
      Array.isArray(user.Role)
        ? user.Role.join(', ')
        : (user.Role ?? 'user') !== value
    )
      return null;
  }

  return children;
}
