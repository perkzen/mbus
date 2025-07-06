const Header = () => {
  return (
    <header
      className={'bg-primary flex w-full flex-row items-center shadow-lg'}
    >
      <img
        alt={'mbus'}
        src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTkYaoql83WuKq0EDreIL2xNkUa0fv9L0u0Ag&s"
        width={64}
        height={64}
      />
      <div className={'ml-4 text-2xl font-bold text-white'}>Mbus</div>
    </header>
  );
};

export default Header;
