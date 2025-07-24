const Header = () => {
  return (
    <header
      className={
        'bg-primary flex w-full flex-row items-center px-6 py-2 md:px-10'
      }
    >
      <img
        alt={'mbus'}
        src="/assets/logo.png"
        className="h-12 w-12 md:h-16 md:w-16"
      />
      <div className="text-background ml-1 text-2xl font-bold tracking-wide md:text-2xl">
        mbus
      </div>
    </header>
  );
};

export default Header;
