const Header = () => {
  return (
    <header
      className={'bg-primary flex w-full flex-row items-center px-10 py-2'}
    >
      <img alt={'mbus'} src="/assets/logo.png" width={64} height={64} />
      <div className={'text-background ml-1 text-2xl font-bold'}>Mbus</div>
    </header>
  );
};

export default Header;
