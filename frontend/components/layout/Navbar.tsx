'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { cn } from '@/lib/utils';

const navItems = [
  { href: '/', label: 'Dashboard' },
  { href: '/entry', label: 'Entry' },
  { href: '/picks', label: 'Picks' },
  { href: '/accumulators', label: 'Accas' },
  { href: '/bets', label: 'Bets' },
  { href: '/performance', label: 'Stats' },
];

export function Navbar() {
  const pathname = usePathname();

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-14 items-center px-6">
        <div className="mr-4 flex">
          <Link href="/" className="mr-6 flex items-center space-x-2">
            <span className="font-bold text-xl">OddsIQ</span>
          </Link>
        </div>
        <nav className="flex items-center space-x-6 text-sm font-medium">
          {navItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                'transition-colors hover:text-foreground/80',
                pathname === item.href
                  ? 'text-foreground font-semibold'
                  : 'text-foreground/60'
              )}
            >
              {item.label}
            </Link>
          ))}
        </nav>
        <div className="ml-auto flex items-center space-x-4">
          <div className="flex items-center space-x-2 text-sm text-muted-foreground">
            <span className="h-2 w-2 rounded-full bg-green-500"></span>
            <span>System Online</span>
          </div>
        </div>
      </div>
    </header>
  );
}
