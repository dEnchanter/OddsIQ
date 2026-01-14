export default function Home() {
  return (
    <main className="min-h-screen p-8">
      <div className="max-w-7xl mx-auto">
        <header className="mb-12">
          <h1 className="text-4xl font-bold mb-2">OddsIQ</h1>
          <p className="text-gray-600 dark:text-gray-400">
            AI-Powered Sports Betting Analytics
          </p>
        </header>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="border border-gray-200 dark:border-gray-800 rounded-lg p-6">
            <h2 className="text-xl font-semibold mb-2">Weekly Picks</h2>
            <p className="text-gray-600 dark:text-gray-400 mb-4">
              View AI-generated betting recommendations for the week
            </p>
            <div className="text-3xl font-bold text-blue-600">0</div>
            <p className="text-sm text-gray-500">Active picks</p>
          </div>

          <div className="border border-gray-200 dark:border-gray-800 rounded-lg p-6">
            <h2 className="text-xl font-semibold mb-2">Performance</h2>
            <p className="text-gray-600 dark:text-gray-400 mb-4">
              Track your betting performance and ROI
            </p>
            <div className="text-3xl font-bold text-green-600">0%</div>
            <p className="text-sm text-gray-500">ROI</p>
          </div>

          <div className="border border-gray-200 dark:border-gray-800 rounded-lg p-6">
            <h2 className="text-xl font-semibold mb-2">Bankroll</h2>
            <p className="text-gray-600 dark:text-gray-400 mb-4">
              Monitor your bankroll and betting history
            </p>
            <div className="text-3xl font-bold text-gray-900 dark:text-gray-100">$0</div>
            <p className="text-sm text-gray-500">Current balance</p>
          </div>
        </div>

        <div className="mt-12 border border-gray-200 dark:border-gray-800 rounded-lg p-6">
          <h2 className="text-2xl font-semibold mb-4">Getting Started</h2>
          <ol className="list-decimal list-inside space-y-2 text-gray-600 dark:text-gray-400">
            <li>Set up your database and run migrations</li>
            <li>Configure API keys for data sources</li>
            <li>Load historical fixture and odds data</li>
            <li>Train the ML model on historical data</li>
            <li>Generate weekly betting picks</li>
          </ol>
        </div>

        <footer className="mt-12 pt-6 border-t border-gray-200 dark:border-gray-800 text-center text-sm text-gray-500">
          <p>OddsIQ MVP - Week 1 Setup Complete</p>
        </footer>
      </div>
    </main>
  );
}
