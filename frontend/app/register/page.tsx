'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/auth';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import Link from 'next/link';
import { UserPlus } from 'lucide-react';

export default function RegisterPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const router = useRouter();
  const { register } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    try {
      await register(username, password);
      router.push('/dashboard');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to register');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-black text-green-400 font-mono">
      <Card className="w-[350px] border-green-500/30 bg-black">
        <CardHeader>
          <CardTitle className="text-2xl font-bold tracking-tighter text-green-400">PARADOX REGISTER</CardTitle>
          <CardDescription className="text-green-500/70">Create a new account to access the dashboard</CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent>
            <div className="space-y-4">
              <div className="space-y-2">
                <label htmlFor="username" className="text-sm font-medium text-green-400">
                  Username
                </label>
                <Input
                  id="username"
                  placeholder="Choose a username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                  className="bg-black border-green-500/30 text-green-400 placeholder:text-green-500/50 focus:border-green-400 focus:ring-green-400/20"
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="password" className="text-sm font-medium text-green-400">
                  Password
                </label>
                <Input
                  id="password"
                  type="password"
                  placeholder="Create a password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  className="bg-black border-green-500/30 text-green-400 placeholder:text-green-500/50 focus:border-green-400 focus:ring-green-400/20"
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="confirm-password" className="text-sm font-medium text-green-400">
                  Confirm Password
                </label>
                <Input
                  id="confirm-password"
                  type="password"
                  placeholder="Confirm your password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  required
                  className="bg-black border-green-500/30 text-green-400 placeholder:text-green-500/50 focus:border-green-400 focus:ring-green-400/20"
                />
              </div>
              {error && (
                <div className="text-sm text-red-500">
                  {error}
                </div>
              )}
            </div>
          </CardContent>
          <CardFooter className="flex flex-col space-y-2">
            <Button 
              type="submit" 
              className="w-full bg-green-900/20 border border-green-500/30 text-green-400 hover:bg-green-900/40"
            >
              <UserPlus className="mr-2 h-4 w-4" />
              Register
            </Button>
            <div className="text-sm text-center text-green-500/70">
              Already have an account?{' '}
              <Link href="/login" className="text-green-400 hover:text-green-300">
                Login
              </Link>
            </div>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
} 