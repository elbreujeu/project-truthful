describe('Env', () => {
    it('should return the correct URL', () => {
      process.env.REACT_APP_API_URL = 'https://example.com/api';
      const { API_URL } = require('../Env');
  
      expect(API_URL).toBe('https://example.com/api');
    });
  });