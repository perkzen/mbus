const Header = () => {
  return (
    <header className={'bg-primary flex w-full flex-row items-center'}>
      <img
        alt={'mbus'}
        src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTkYaoql83WuKq0EDreIL2xNkUa0fv9L0u0Ag&s"
        width={64}
        height={64}
      />
      <div className={'text-background ml-4 text-2xl font-bold'}>Mbus</div>
    </header>
  );
};

export default Header;
