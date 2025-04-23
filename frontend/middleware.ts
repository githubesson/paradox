import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';


const publicPaths = ['/login', '/register'];

export function middleware(request: NextRequest) {
  const token = request.cookies.get('token');
  const { pathname } = request.nextUrl;

  
  if (publicPaths.includes(pathname)) {
    
    if (token) {
      return NextResponse.redirect(new URL('/dashboard', request.url));
    }
    return NextResponse.next();
  }

  
  if (!token) {
    
    const loginUrl = new URL('/login', request.url);
    loginUrl.searchParams.set('from', pathname);
    return NextResponse.redirect(loginUrl);
  }

  return NextResponse.next();
}

export const config = {
  
  matcher: [
    /*
     * Match all paths except:
     * 1. /api (API routes)
     * 2. /_next (Next.js internals)
     * 3. /static (static files)
     * 4. .*\\..*$ (files with extensions)
     */
    '/((?!api|_next|static|.*\\..*$).*)',
  ],
}; 